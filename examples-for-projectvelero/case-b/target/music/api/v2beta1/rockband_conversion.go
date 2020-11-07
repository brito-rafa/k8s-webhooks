package v2beta1

import (
	v2beta2 "music/api/v2beta2"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (

	// this is the annotation key to keep drummer value if converted from v2beta2 to v2veta1
	drummerAnnotation = "rockbands.v2beta2.music.example.io/drummer"
	// default drummer string to be used when converting from v1 to v2beta2
	defaultValueDrummerConverter = "Converted to v2beta2"
)

// ConvertTo converts this RockBand v2beta1 to the Hub version (v2beta2)
func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v2beta2.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v2beta1-to-v2beta2")

	rockbandlog.Info("ConvertTo v2beta2 from v2beta1", "name", src.Name, "namespace", src.Namespace, "lead guitar", src.Spec.LeadGuitar)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	annotations := src.GetAnnotations()

	if annotations != nil && annotations[drummerAnnotation] != "" {
		dst.Spec.Drummer = annotations[drummerAnnotation]
		rockbandlog.Info("ConvertTo v2beta2 from v2beta1 - found annotations", "name", dst.Name, "namespace", dst.Namespace, "drummer", dst.Spec.Drummer)
	} else {
		// Setting a default string as leadSinger
		dst.Spec.Drummer = defaultValueDrummerConverter
		rockbandlog.Info("ConvertTo v2beta2 from v2beta1 - no annotations, using the default", "name", dst.Name, "namespace", dst.Namespace, "drummer", dst.Spec.Drummer)
	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger
	dst.Spec.LeadGuitar = src.Spec.LeadGuitar

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}

// ConvertFrom converts from the Hub version (v2beta2) to the version v2beta1
func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v2beta2.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v2beta2-to-v2beta1")

	rockbandlog.Info("ConvertFrom v2beta2 to v2beta1", "name", src.Name, "namespace", src.Namespace)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.Drummer != defaultValueDrummerConverter {
		annotations := dst.GetAnnotations()

		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations[drummerAnnotation] = src.Spec.Drummer

		dst.SetAnnotations(annotations)

		rockbandlog.Info("ConvertTo v2beta2 from v2beta1 - set annotations", "name", src.Name, "namespace", src.Namespace)

	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger
	dst.Spec.LeadGuitar = src.Spec.LeadGuitar

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}
