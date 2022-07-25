package mpeg

// https://wiki.multimedia.cx/index.php/MPEG-4_Audio

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var mpegAscMpegUSACCFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.MPEG_ASC,
		Description: "MPEG-4 Audio Specific Config",
		DecodeFn:    ascDecoder,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_USACC}, Group: &mpegAscMpegUSACCFormat},
		},
	})
}

var frequencyIndexHzMap = scalar.UToSymU{
	0x0: 96000,
	0x1: 88200,
	0x2: 64000,
	0x3: 48000,
	0x4: 44100,
	0x5: 32000,
	0x6: 24000,
	0x7: 22050,
	0x8: 16000,
	0x9: 12000,
	0xa: 11025,
	0xb: 8000,
	0xc: 7350,
}

var channelConfigurationNames = scalar.UToDescription{
	0: "defined in AOT Specifc Config",
	1: "front-center",
	2: "front-left, front-right",
	3: "front-center, front-left, front-right",
	4: "front-center, front-left, front-right, back-center",
	5: "front-center, front-left, front-right, back-left, back-right",
	6: "front-center, front-left, front-right, back-left, back-right, LFE-channel",
	7: "front-center, front-left, front-right, side-left, side-right, back-left, back-right, LFE-channel",
}

func ascDecoder(d *decode.D, _ any) any {
	objectType := d.FieldUFn("object_type", decodeEscapeValueCarryFn(5, 6, 0), format.MPEGAudioObjectTypeNames)
	d.FieldUFn("sampling_frequency", decodeEscapeValueAbsFn(4, 24, 0), frequencyIndexHzMap)
	d.FieldU4("channel_configuration", channelConfigurationNames)
	// TODO: GASpecificConfig etc

	switch objectType {
	case format.MPEGAudioObjectTypeUSAC:
		d.FieldFormat("usac_config", mpegAscMpegUSACCFormat, nil)
	default:
		// TODO: GASpecificConfig etc
		d.FieldRawLen("var_aot_or_byte_align", d.BitsLeft())
	}

	return format.MPEGASCOut{ObjectType: int(objectType)}
}
