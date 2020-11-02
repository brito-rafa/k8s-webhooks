package v1alpha1

import (
	v1 "music/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

var (
	leadSingerAnnotation = "rockband.music.example.io/lead-singer"
)

func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.RockBand)
	dst.Spec.LeadSinger = "TBD Converter"
	return nil
}

func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.RockBand)
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	// Retaining the LeadSinger as an annotation
	dst.Annotations[leadSingerAnnotation] = src.Spec.LeadSinger
	return nil
}
