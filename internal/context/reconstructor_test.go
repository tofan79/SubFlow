package context

import "testing"

func TestReconstruct_MergeFragmented(t *testing.T) {
	segments := []Segment{
		{StartMS: 0, EndMS: 1000, Text: "Hello"},
		{StartMS: 1200, EndMS: 2000, Text: "world."},
	}

	out, scenes := Reconstruct(segments, DefaultReconstructConfig())
	if len(out) != 1 {
		t.Fatalf("expected 1 segment, got %d: %#v", len(out), out)
	}
	if out[0].Text != "Hello world." {
		t.Fatalf("merged text mismatch: %#v", out[0])
	}
	if !out[0].IsMerged {
		t.Fatalf("expected merged segment: %#v", out[0])
	}
	if out[0].StartMS != 0 || out[0].EndMS != 2000 {
		t.Fatalf("merged timings mismatch: %#v", out[0])
	}
	if len(scenes) != 1 {
		t.Fatalf("expected 1 scene, got %d: %#v", len(scenes), scenes)
	}
	if len(scenes[0].Segments) != 2 || scenes[0].Segments[0] != 1 || scenes[0].Segments[1] != 2 {
		t.Fatalf("scene segment indices mismatch: %#v", scenes[0])
	}
}

func TestReconstruct_SceneDetection(t *testing.T) {
	segments := []Segment{
		{StartMS: 0, EndMS: 1000, Text: "One."},
		{StartMS: 5000, EndMS: 6000, Text: "Two."},
		{StartMS: 10000, EndMS: 11000, Text: "Three."},
	}

	out, scenes := Reconstruct(segments, DefaultReconstructConfig())
	if len(out) != 3 {
		t.Fatalf("expected 3 segments, got %d: %#v", len(out), out)
	}
	if len(scenes) != 3 {
		t.Fatalf("expected 3 scenes, got %d: %#v", len(scenes), scenes)
	}
	for i, seg := range out {
		if seg.SceneID != i+1 {
			t.Fatalf("segment %d scene mismatch: %#v", i, seg)
		}
	}
	for i, scene := range scenes {
		if scene.SceneID != i+1 {
			t.Fatalf("scene %d id mismatch: %#v", i, scene)
		}
		if len(scene.Segments) != 1 || scene.Segments[0] != i+1 {
			t.Fatalf("scene %d segments mismatch: %#v", i, scene)
		}
	}
}

func TestReconstruct_NoMergeComplete(t *testing.T) {
	segments := []Segment{
		{StartMS: 0, EndMS: 1000, Text: "Hello."},
		{StartMS: 1100, EndMS: 2000, Text: "World."},
	}

	out, scenes := Reconstruct(segments, DefaultReconstructConfig())
	if len(out) != 2 {
		t.Fatalf("expected 2 segments, got %d: %#v", len(out), out)
	}
	if out[0].Text != "Hello." || out[1].Text != "World." {
		t.Fatalf("sentences should remain unchanged: %#v", out)
	}
	if out[0].IsMerged || out[1].IsMerged {
		t.Fatalf("complete sentences should not merge: %#v", out)
	}
	if len(scenes) != 1 {
		t.Fatalf("expected 1 scene, got %d: %#v", len(scenes), scenes)
	}
}
