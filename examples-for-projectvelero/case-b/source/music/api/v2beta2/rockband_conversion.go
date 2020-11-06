package v2beta2

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

	// this is the annotation key to keep drummer value if converted from v2beta2 to v1
	drummerAnnotation = "rockbands.v2beta2.music.example.io/drummer"
	// default drummer string to be used when converting from v1 to v2beta2
	defaultValueDrummerConverter = "Converted from v1"
)

// ConvertTo converts this RockBand v2beta2 to the Hub version (v1)
func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v2beta2-to-v1")

	rockbandlog.Info("ConvertTo v1 from v2beta2", "name", src.Name, "namespace", src.Namespace, "lead guitar", src.Spec.LeadGuitar, "drummer", src.Spec.Drummer)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	annotations := dst.GetAnnotations()

	if annotations == nil {
		annotations = make(map[string]string)
	}

	if src.Spec.LeadGuitar != defaultValueLeadGuitarConverter {
		annotations[leadGuitarAnnotation] = src.Spec.LeadGuitar
		rockbandlog.Info("ConvertTo v1 from v2beta2 - set annotations on leadGuitar", "name", src.Name, "namespace", src.Namespace, "leadGuitar", src.Spec.LeadGuitar)
	}

	if src.Spec.Drummer != defaultValueDrummerConverter {
		annotations[drummerAnnotation] = src.Spec.Drummer
		rockbandlog.Info("ConvertTo v1 from v2beta2 - set annotations on drummer", "name", src.Name, "namespace", src.Namespace, "drummer", src.Spec.Drummer)
	}

	dst.SetAnnotations(annotations)

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}

// ConvertFrom converts from the Hub version (v1) to the version v2beta2
func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v1-to-v2beta2")

	rockbandlog.Info("ConvertFrom v1 to v2beta2", "name", src.Name, "namespace", src.Namespace)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	annotations := src.GetAnnotations()

	if annotations != nil && annotations[leadGuitarAnnotation] != "" {
		dst.Spec.LeadGuitar = annotations[leadGuitarAnnotation]
		rockbandlog.Info("ConvertFrom v1 to v2beta2 - found leadGuitar annotation", "name", dst.Name, "namespace", dst.Namespace, "lead guitar", dst.Spec.LeadGuitar)
	} else {
		// Setting a default string as leadSinger
		dst.Spec.LeadGuitar = defaultValueLeadGuitarConverter
		rockbandlog.Info("ConvertFrom v1 to v2beta2 - no leadGuitar annotations, using the default", "name", dst.Name, "namespace", dst.Namespace, "lead guitar", dst.Spec.LeadGuitar)
	}

	if annotations != nil && annotations[drummerAnnotation] != "" {
		dst.Spec.Drummer = annotations[drummerAnnotation]
		rockbandlog.Info("ConvertFrom v1 to v2beta2 - found drummer annotations", "name", dst.Name, "namespace", dst.Namespace, "drummer", dst.Spec.Drummer)
	} else {
		// Setting a default string as leadSinger
		dst.Spec.Drummer = defaultValueDrummerConverter
		rockbandlog.Info("ConvertFrom v1 to v2beta2 - no drummer annotations, using the default", "name", dst.Name, "namespace", dst.Namespace, "drummer", dst.Spec.Drummer)
	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}
