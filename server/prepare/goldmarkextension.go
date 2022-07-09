package prepare

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Segment AST Node

type Segment struct {
	ast.BaseBlock
	Segments text.Segments
}

func NewSegment() *Segment {
	return &Segment{
		BaseBlock: ast.BaseBlock{},
	}
}

func (n *Segment) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindSegment = ast.NewNodeKind("Segment")

func (n *Segment) Kind() ast.NodeKind {
	return KindSegment
}

// Segment renderer

type SegmentRenderer struct{}

func NewSegmentRenderer() renderer.NodeRenderer {
	r := &SegmentRenderer{}
	return r
}

func (r *SegmentRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindSegment, r.renderSegment)
}

func (r *SegmentRenderer) renderSegment(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<span class=\"segment\">")
	} else {
		_, _ = w.WriteString("</span>\n")
	}
	return ast.WalkContinue, nil
}

// HackParagraph AST Node

type HackParagraph struct {
	ast.BaseBlock
	Segments text.Segments
}

func NewHackParagraph() *HackParagraph {
	return &HackParagraph{
		BaseBlock: ast.BaseBlock{},
	}
}

func (n *HackParagraph) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindHackParagraph = ast.NewNodeKind("HackParagraph")

func (n *HackParagraph) Kind() ast.NodeKind {
	return KindHackParagraph
}

// HackParagraph renderer

type HackParagraphRenderer struct{}

func NewHackParagraphRenderer() renderer.NodeRenderer {
	r := &HackParagraphRenderer{}
	return r
}

func (r *HackParagraphRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindHackParagraph, r.renderHackParagraph)
}

func (r *HackParagraphRenderer) renderHackParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<p>")
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return ast.WalkContinue, nil
}

// ParagraphTransformer that replaces every Paragraph with a HackParagraph whose children are Segments

type segmentParagraphTransformer struct {
}

var SegmentParagraphTransformer = &segmentParagraphTransformer{}

func (p *segmentParagraphTransformer) Transform(node *ast.Paragraph, reader text.Reader, pc parser.Context) {
	newPar := NewHackParagraph()

	block := text.NewBlockReader(reader.Source(), node.Lines())
	segment := NewSegment()
	for {
		line, lastSeg := block.PeekLine()

		if line == nil || util.IsBlank(line) {
			if segment.Lines().Len() > 0 {
				newPar.AppendChild(newPar, segment)
			}
			break
		}

		for _, b := range line {
			block.Advance(1)
			if b != '.' && b != '?' && b != '!' {
				continue
			}
			if block.Peek() != ' ' {
				// Don't split if the tentative end-of-segment marker is not followed by a space.
				continue
			}
			lastSeg = lastSeg.TrimLeftSpace(reader.Source())
			_, currSeg := block.Position()
			newSeg := lastSeg.Between(currSeg)
			// Don't split if the resulting segment would be very short
			if newSeg.Len() < 16 && segment.Lines().Len() == 0 {
				continue
			}
			// Split here
			segment.Lines().Append(newSeg)
			lastSeg = currSeg
			newPar.AppendChild(newPar, segment)
			segment = NewSegment()
		}

		// Commit the line
		if lastSeg.Len() > 1 {
			lastSeg = lastSeg.TrimLeftSpace(reader.Source())
			segment.Lines().Append(lastSeg)
		}
	}

	// Replace the nasty old Paragraph with the cool HackParagraph
	node.Parent().AppendChild(node.Parent(), newPar)
	node.Parent().RemoveChild(node.Parent(), node)
}

type coditraGoldmark struct{}

var CoditraGoldmark = &coditraGoldmark{}

func NewCoditraGoldmark() goldmark.Extender {
	return &coditraGoldmark{}
}

func (e *coditraGoldmark) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithParagraphTransformers(
			util.Prioritized(SegmentParagraphTransformer, 50),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHackParagraphRenderer(), 50),
			util.Prioritized(NewSegmentRenderer(), 60),
		),
	)
}
