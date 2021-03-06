package mp4

import "io"

// Media Information Box (minf - mandatory)
//
// Contained in : Media Box (mdia)
//
// Status: partially decoded
type MinfBox struct {
	//Nmhd *NmhdBox `json:"nmhd,"`
	Vmhd  *VmhdBox `json:"vmhd,"`
	Smhd  *SmhdBox `json:"smhd,"`
	Stbl  *StblBox `json:"stbl,"`
	Dinf  *DinfBox `json:"dinf,"`
	Hdlr  *HdlrBox `json:"hdlr,omitempty"` // ISO IEC 14496-12, does not contain hdlr in minf
	Boxes []Box    `json:",omitempty"`
}

func DecodeMinf(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MinfBox{}
	for _, b := range l {
		switch b.Type() {
		case "vmhd":
			m.Vmhd = b.(*VmhdBox)
		case "smhd":
			m.Smhd = b.(*SmhdBox)
		case "stbl":
			m.Stbl = b.(*StblBox)
		case "dinf":
			m.Dinf = b.(*DinfBox)
		case "hdlr":
			m.Hdlr = b.(*HdlrBox)
		default:
			m.Boxes = append(m.Boxes, b.Box())
		}
	}
	return m, nil
}

func (b *MinfBox) Box() Box {
	return b
}

func (b *MinfBox) Type() string {
	return "minf"
}

func (b *MinfBox) Size() int {
	sz := 0
	if b.Vmhd != nil {
		sz += b.Vmhd.Size()
	}
	if b.Smhd != nil {
		sz += b.Smhd.Size()
	}
	sz += b.Stbl.Size()
	if b.Dinf != nil {
		sz += b.Dinf.Size()
	}
	if b.Hdlr != nil {
		sz += b.Hdlr.Size()
	}
	return sz + BoxHeaderSize
}

func (b *MinfBox) Dump() {
	b.Stbl.Dump()
}

func (b *MinfBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	if b.Vmhd != nil {
		err = b.Vmhd.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Smhd != nil {
		err = b.Smhd.Encode(w)
		if err != nil {
			return err
		}
	}
	err = b.Dinf.Encode(w)
	if err != nil {
		return err
	}
	err = b.Stbl.Encode(w)
	if err != nil {
		return err
	}
	if b.Hdlr != nil {
		return b.Hdlr.Encode(w)
	}
	return nil
}
