package v1alpha1

import (
	v1 "music/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

var (
	// this is the annotation key to keep leadSinger value if converted from v1 to v1alpha1
	leadSingerAnnotation = "rockbands.v1.music.example.io/leadSinger"
	// default leadSinger string to be used when converting from v1alpha1 to v1
	defaultValueLeadSingerConverter = "Converted from v1alpha1"
)

// ConvertTo converts this RockBand v1alpha1 to the Hub version (v1)
func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.RockBand)

	// Checking if there is a leadSinger as annotation
	// If so, use this value to set the leadSinger

	// This condition will only be met if
	// you are creating a v1alpha1 object with
	// v1 leadSinger field as annotations
	// see sample yaml under music/config/samples

	annotations := src.GetAnnotations()

	if annotations != nil && annotations[leadSingerAnnotation] != "" {
		dst.Spec.LeadSinger = annotations[leadSingerAnnotation]
	} else {
		// Setting a default string as leadSinger
		dst.Spec.LeadSinger = defaultValueLeadSingerConverter
	}

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version valpha1
// This is what the Velero backup will see: v1alpha version
func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.RockBand)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Retaining the v1 LeadSinger values as an annotation
	// Saving fields as annotation is a great way
	// to keep information back and forth between legacy and modern API Groups

	// if the leadSinger is already is set as the default value from v1alpha1 (see ConvertTo)
	// do not bother to create an annotation

	if src.Spec.LeadSinger != defaultValueLeadSingerConverter {
		annotations := dst.GetAnnotations()

		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations[leadSingerAnnotation] = src.Spec.LeadSinger

		dst.SetAnnotations(annotations)

	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}
