package v2beta2

import (
	v2 "music/api/v2"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (

	// this is the annotation key to keep bass value if converted from v2 to v2beta1
	bassAnnotation = "rockbands.v2.music.example.io/bass"
	// default bass string to be used when converting from v2beta1 to v2
	defaultValueBassConverter = "Converted from v2beta1 to v2"
)

// ConvertTo converts this RockBand v2beta2 to the Hub version (v2)
func (src *RockBand) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v2.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v2beta2-to-v2")

	rockbandlog.Info("ConvertTo v2 from v2beta2", "name", src.Name, "namespace", src.Namespace, "lead guitar", src.Spec.LeadGuitar)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	annotations := src.GetAnnotations()

	if annotations != nil && annotations[bassAnnotation] != "" {
		dst.Spec.Bass = annotations[bassAnnotation]
		rockbandlog.Info("ConvertTo v2 from v2beta2 - found annotations", "name", dst.Name, "namespace", dst.Namespace, "bass", dst.Spec.Bass)
	} else {
		// Setting a default string as bass
		dst.Spec.Bass = defaultValueBassConverter
		rockbandlog.Info("ConvertTo v2 from v2beta2 - no annotations, using the default", "name", dst.Name, "namespace", dst.Namespace, "bass", dst.Spec.Bass)
	}

	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger
	dst.Spec.LeadGuitar = src.Spec.LeadGuitar
	dst.Spec.Drummer = src.Spec.Drummer

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}

// ConvertFrom converts from the Hub version (v2) to the version v2beta2
func (dst *RockBand) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v2.RockBand)

	rockbandlog = logf.Log.WithName("rockband-convert-v2-to-v2beta2")

	rockbandlog.Info("ConvertFrom v2 to v2beta2", "name", src.Name, "namespace", src.Namespace)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	annotations := dst.GetAnnotations()

	if annotations == nil {
		annotations = make(map[string]string)
	}

	if src.Spec.Bass != defaultValueBassConverter {
		annotations[bassAnnotation] = src.Spec.Bass
		rockbandlog.Info("ConvertTo v2 from v2beta2 - set annotations", "name", src.Name, "namespace", src.Namespace, "bass", src.Spec.Bass)
	}

	dst.SetAnnotations(annotations)
	// Other Spec
	dst.Spec.NumberComponents = src.Spec.NumberComponents
	dst.Spec.Genre = src.Spec.Genre
	dst.Spec.LeadSinger = src.Spec.LeadSinger
	dst.Spec.LeadGuitar = src.Spec.LeadGuitar
	dst.Spec.Drummer = src.Spec.Drummer

	// Status
	dst.Status.LastPlayed = src.Status.LastPlayed

	return nil
}
