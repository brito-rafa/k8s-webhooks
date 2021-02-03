package v2beta1

import (
	v1 "music/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	// this is the annotation key to keep leadGuitar value if converted from v2beta1 to v1
	leadGuitarAnnotation = "rockbands.v2beta1.music.example.io/leadGuitar"
	// default leadSinger string to be used when converting from v1 to v2beta1
	defaultValueLeadGuitarConverter = "Converted from v1"
)

// ConvertTo converts this RockBand v2beta1 to the Hub version (v1)
func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v2beta1-to-v1")

	rockbandlog.Info("ConvertTo v1 from v2beta1", "name", src.Name, "namespace", src.Namespace, "lead guitar", src.Spec.LeadGuitar)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.LeadGuitar != defaultValueLeadGuitarConverter {
		annotations := dst.GetAnnotations()

		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations[leadGuitarAnnotation] = src.Spec.LeadGuitar

		dst.SetAnnotations(annotations)

		rockbandlog.Info("ConvertTo v1 from v2beta1 - set annotations", "name", src.Name, "namespace", src.Namespace)

	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}

// ConvertFrom converts from the Hub version (v1) to the version v2beta1
func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v1-to-v2beta1")

	rockbandlog.Info("ConvertFrom v1 to v2beta1 - lead guitar is undef at this point", "name", src.Name, "namespace", src.Namespace)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	annotations := src.GetAnnotations()

	if annotations != nil && annotations[leadGuitarAnnotation] != "" {
		dst.Spec.LeadGuitar = annotations[leadGuitarAnnotation]
		rockbandlog.Info("ConvertFrom v1 to v2beta1 - found annotations", "name", dst.Name, "namespace", dst.Namespace, "lead guitar", dst.Spec.LeadGuitar)
	} else {
		// Setting a default string as leadSinger
		dst.Spec.LeadGuitar = defaultValueLeadGuitarConverter
		rockbandlog.Info("ConvertFrom v1 to v2beta1 - no annotations, using the default", "name", dst.Name, "namespace", dst.Namespace, "lead guitar", dst.Spec.LeadGuitar)
	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}
