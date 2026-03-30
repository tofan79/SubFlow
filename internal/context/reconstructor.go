package context

import "unicode/utf8"

type Segment struct {
	Index    int
	StartMS  int64
	EndMS    int64
	Text     string
	Speaker  string
	SceneID  int
	IsMerged bool
}

type SceneMetadata struct {
	SceneID  int
	StartMS  int64
	EndMS    int64
	Segments []int
}

type ReconstructConfig struct {
	MaxGapForMergeMS int64
	SceneGapMS       int64
}

func DefaultReconstructConfig() ReconstructConfig {
	return ReconstructConfig{
		MaxGapForMergeMS: 500,
		SceneGapMS:       3000,
	}
}

func Reconstruct(segments []Segment, cfg ReconstructConfig) ([]Segment, []SceneMetadata) {
	cfg = normalizeReconstructConfig(cfg)
	if len(segments) == 0 {
		return []Segment{}, []SceneMetadata{}
	}

	out := make([]Segment, 0, len(segments))
	scenes := make([]SceneMetadata, 0, len(segments))
	currentSceneIdx := -1

	for i, seg := range segments {
		segmentIndex := seg.Index
		if segmentIndex <= 0 {
			segmentIndex = i + 1
		}

		if len(out) > 0 && ShouldMerge(out[len(out)-1], seg, cfg.MaxGapForMergeMS) {
			last := &out[len(out)-1]
			last.EndMS = seg.EndMS
			last.Text = mergeText(last.Text, seg.Text)
			if last.Speaker == "" {
				last.Speaker = seg.Speaker
			}
			last.IsMerged = true
			if currentSceneIdx >= 0 {
				scenes[currentSceneIdx].EndMS = last.EndMS
				scenes[currentSceneIdx].Segments = append(scenes[currentSceneIdx].Segments, segmentIndex)
			}
			continue
		}

		seg.Index = len(out) + 1
		out = append(out, seg)

		if len(out) == 1 {
			scenes = append(scenes, SceneMetadata{
				SceneID:  1,
				StartMS:  seg.StartMS,
				EndMS:    seg.EndMS,
				Segments: []int{segmentIndex},
			})
			currentSceneIdx = 0
			out[0].SceneID = 1
			continue
		}

		prev := out[len(out)-2]
		gap := seg.StartMS - prev.EndMS
		if gap > cfg.SceneGapMS {
			scenes = append(scenes, SceneMetadata{
				SceneID:  len(scenes) + 1,
				StartMS:  seg.StartMS,
				EndMS:    seg.EndMS,
				Segments: []int{segmentIndex},
			})
			currentSceneIdx = len(scenes) - 1
		} else {
			scenes[currentSceneIdx].EndMS = seg.EndMS
			scenes[currentSceneIdx].Segments = append(scenes[currentSceneIdx].Segments, segmentIndex)
		}

		out[len(out)-1].SceneID = scenes[currentSceneIdx].SceneID
	}

	return out, scenes
}

func ShouldMerge(a, b Segment, maxGapMS int64) bool {
	if maxGapMS <= 0 {
		return false
	}
	if b.StartMS-a.EndMS >= maxGapMS {
		return false
	}
	return !endsWithSentencePunctuation(a.Text)
}

func normalizeReconstructConfig(cfg ReconstructConfig) ReconstructConfig {
	def := DefaultReconstructConfig()
	if cfg.MaxGapForMergeMS <= 0 {
		cfg.MaxGapForMergeMS = def.MaxGapForMergeMS
	}
	if cfg.SceneGapMS <= 0 {
		cfg.SceneGapMS = def.SceneGapMS
	}
	return cfg
}

func mergeText(a, b string) string {
	switch {
	case a == "":
		return b
	case b == "":
		return a
	default:
		return a + " " + b
	}
}

func endsWithSentencePunctuation(text string) bool {
	if text == "" {
		return false
	}
	r, _ := utf8.DecodeLastRuneInString(text)
	switch r {
	case '.', '!', '?':
		return true
	default:
		return false
	}
}
