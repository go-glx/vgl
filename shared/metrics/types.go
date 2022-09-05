package metrics

type Segment = string

const (
	SegmentPlClearBuffers        Segment = "pl.clear.buff"
	SegmentPlUpdateGlobalUniform Segment = "pl.upd.ubo"
	SegmentPlUpdateSSBO          Segment = "pl.upd.ssbo"
	SegmentPlUpdateIndexes       Segment = "pl.upd.ind"
	SegmentPlUpdateVertexes      Segment = "pl.upd.vert"
	SegmentPlCreatePipeline      Segment = "pl.create.pipe"
	SegmentPlBindPipeline        Segment = "pl.bind.pipe"
	SegmentPlBindIndexes         Segment = "pl.bind.ind"
	SegmentPlBindVertex          Segment = "pl.bind.vert"
	SegmentPlBindUniforms        Segment = "pl.bind.uniform"
	SegmentPlDraw                Segment = "pl.draw"
)
