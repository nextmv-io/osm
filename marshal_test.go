package osm

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestChangeCompare(t *testing.T) {
	data := readFile(t, "testdata/changeset_38162206.osc")

	c1 := &Change{}
	err := xml.Unmarshal(data, &c1)
	if err != nil {
		t.Errorf("unable to unmarshal: %v", err)
	}

	c2 := &Change{}
	err = xml.Unmarshal(data, &c2)
	if err != nil {
		t.Errorf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("changes are not equal")
	}
}

func TestProtobufNode(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162210.osc")
	n1 := c.Create.Nodes[12]

	// verify it's a good test
	if len(n1.Tags) == 0 {
		t.Fatalf("test should have some tags")
	}

	ss := &stringSet{}
	pbnode := marshalNode(n1, ss, true)

	n2, err := unmarshalNode(pbnode, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(n1, n2) {
		t.Errorf("nodes are not equal")
		t.Logf("%+v", n1)
		t.Logf("%+v", n2)
	}
}

func TestProtobufNodeRoundoff(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162210.osc")
	n1 := c.Create.Nodes[194]

	ss := &stringSet{}
	pbnode := marshalNode(n1, ss, true)

	n2, err := unmarshalNode(pbnode, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(n1, n2) {
		t.Errorf("nodes are not equal")
		t.Logf("%+v", n1)
		t.Logf("%+v", n2)
	}
}

func TestProtobufNodes(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162210.osc")
	ns1 := c.Create.Nodes

	ss := &stringSet{}
	pbnodes := marshalNodes(ns1, ss, true)

	ns2, err := unmarshalNodes(pbnodes, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if len(ns1) != len(ns2) {
		t.Fatalf("different number of nodes: %d != %d", len(ns1), len(ns2))
	}

	for i := range ns1 {
		ns1[i].XMLName = xml.Name{}
		ns2[i].XMLName = xml.Name{}
		if !reflect.DeepEqual(ns1[i], ns2[i]) {
			t.Errorf("nodes %d are not equal", i)
			t.Logf("%+v", ns1[i])
			t.Logf("%+v", ns2[i])
		}
	}

	// nodes with no tags
	for _, n := range ns1 {
		n.Tags = nil
	}

	ss = &stringSet{}
	pbnodes = marshalNodes(ns1, ss, true)

	ns2, err = unmarshalNodes(pbnodes, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if len(ns1) != len(ns2) {
		t.Fatalf("different number of nodes: %d != %d", len(ns1), len(ns2))
	}

	for i := range ns1 {
		ns1[i].XMLName = xml.Name{}
		ns2[i].XMLName = xml.Name{}
		if !reflect.DeepEqual(ns1[i], ns2[i]) {
			t.Errorf("nodes %d are not equal", i)
			t.Logf("%+v", ns1[i])
			t.Logf("%+v", ns2[i])
		}
	}
}

func TestProtobufWay(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162210.osc")
	w1 := c.Create.Ways[5]

	// verify it's a good test
	if len(w1.Tags) == 0 {
		t.Fatalf("test should have some tags")
	}

	ss := &stringSet{}
	pbway := marshalWay(w1, ss, true)

	w2, err := unmarshalWay(pbway, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(w1, w2) {
		t.Errorf("ways are not equal")
		t.Logf("%+v", w1)
		t.Logf("%+v", w2)
	}
}

func TestProtobufMinorWay(t *testing.T) {
	o := loadOSM(t, "testdata/minor-way.osm")
	w1 := o.Ways[0]

	ss := &stringSet{}
	pbway := marshalWay(w1, ss, true)

	w2, err := unmarshalWay(pbway, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(w1, w2) {
		t.Errorf("ways are not equal")
		t.Logf("%+v", w1)
		t.Logf("%+v", w2)
	}

	// with no minor nodes
	w1.Minors[0].MinorNodes = nil
	w1.Minors[1].MinorNodes = nil

	ss = &stringSet{}
	pbway = marshalWay(w1, ss, true)

	w2, err = unmarshalWay(pbway, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(w1, w2) {
		t.Errorf("ways are not equal")
		t.Logf("%+v", w1)
		t.Logf("%+v", w2)
	}

	// with no minor version
	w1.Minors = nil
	w1.Minors = nil

	ss = &stringSet{}
	pbway = marshalWay(w1, ss, true)

	w2, err = unmarshalWay(pbway, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(w1, w2) {
		t.Errorf("ways are not equal")
		t.Logf("%+v", w1)
		t.Logf("%+v", w2)
	}
}

func TestProtobufMinorRelation(t *testing.T) {
	o := loadOSM(t, "testdata/minor-relation.osm")
	r1 := o.Relations[0]

	ss := &stringSet{}
	pbrelation := marshalRelation(r1, ss, true)

	r2, err := unmarshalRelation(pbrelation, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("relations are not equal")
		t.Logf("%+v", r1)
		t.Logf("%+v", r2)
	}

	// with no minor members
	r1.Minors[0].MinorMembers = nil
	r1.Minors[1].MinorMembers = nil

	ss = &stringSet{}
	pbrelation = marshalRelation(r1, ss, true)

	r2, err = unmarshalRelation(pbrelation, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("relations are not equal")
		t.Logf("%+v", r1)
		t.Logf("%+v", r2)
	}

	// with no minor version
	r1.Minors = nil
	r1.Minors = nil

	ss = &stringSet{}
	pbrelation = marshalRelation(r1, ss, true)

	r2, err = unmarshalRelation(pbrelation, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("relations are not equal")
		t.Logf("%+v", r1)
		t.Logf("%+v", r2)
	}
}

func TestProtobufRelation(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162206.osc")
	r1 := c.Create.Relations[0]

	// verify it's a good test
	if len(r1.Tags) == 0 {
		t.Fatalf("test should have some tags")
	}

	ss := &stringSet{}
	pbrelation := marshalRelation(r1, ss, true)

	r2, err := unmarshalRelation(pbrelation, ss.Strings(), nil)
	if err != nil {
		t.Fatalf("unable to unmarshal: %v", err)
	}

	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("relations are not equal")
		t.Logf("%+v", r1)
		t.Logf("%+v", r2)
	}
}

func BenchmarkMarshalXML(b *testing.B) {
	filename := "testdata/changeset_38162206.osc"
	data := readFile(b, filename)

	c := &Change{}
	err := xml.Unmarshal(data, c)
	if err != nil {
		b.Fatalf("unable to unmarshal: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := xml.Marshal(c)
		if err != nil {
			b.Fatalf("unable to marshal: %v", err)
		}
	}
}

func BenchmarkMarshalProto(b *testing.B) {
	cs := &Changeset{
		ID:     38162206,
		UserID: 2744209,
		User:   "grah735",
		Change: loadChange(b, "testdata/changeset_38162206.osc"),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := cs.Marshal()
		if err != nil {
			b.Fatalf("unable to marshal: %v", err)
		}
	}
}

func BenchmarkMarshalMinorWay(b *testing.B) {
	o := loadOSM(b, "testdata/minor-way.osm")
	ways := o.Ways

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data, _ := ways.Marshal()
		UnmarshalWays(data)
	}
}

func BenchmarkMarshalMinorRelation(b *testing.B) {
	o := loadOSM(b, "testdata/minor-relation.osm")
	relations := o.Relations

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data, _ := relations.Marshal()
		UnmarshalRelations(data)
	}
}

func BenchmarkMarshalProtoGZIP(b *testing.B) {
	cs := &Changeset{
		ID:     38162206,
		UserID: 2744209,
		User:   "grah735",
		Change: loadChange(b, "testdata/changeset_38162206.osc"),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		data, err := cs.Marshal()
		if err != nil {
			b.Fatalf("unable to marshal: %v", err)
		}

		w, _ := gzip.NewWriterLevel(&bytes.Buffer{}, gzip.BestCompression)
		_, err = w.Write(data)
		if err != nil {
			b.Fatalf("unable to write: %v", err)
		}

		err = w.Close()
		if err != nil {
			b.Fatalf("unable to close: %v", err)
		}
	}
}

func BenchmarkUnmarshalXML(b *testing.B) {
	filename := "testdata/changeset_38162206.osc"
	data := readFile(b, filename)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		c := &Change{}
		err := xml.Unmarshal(data, c)
		if err != nil {
			b.Fatalf("unable to unmarshal: %v", err)
		}
	}
}

func BenchmarkUnmarshalProto(b *testing.B) {
	cs := &Changeset{
		ID:     38162206,
		UserID: 2744209,
		User:   "grah735",
		Change: loadChange(b, "testdata/changeset_38162206.osc"),
	}

	data, err := cs.Marshal()
	if err != nil {
		b.Fatalf("unable to marshal: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := UnmarshalChangeset(data)
		if err != nil {
			b.Fatalf("unable to unmarshal: %v", err)
		}
	}
}
func BenchmarkUnmarshalProtoGZIP(b *testing.B) {
	cs := &Changeset{
		ID:     38162206,
		UserID: 2744209,
		User:   "grah735",
		Change: loadChange(b, "testdata/changeset_38162206.osc"),
	}

	data, err := cs.Marshal()
	if err != nil {
		b.Fatalf("unable to marshal: %v", err)
	}

	buff := &bytes.Buffer{}
	w, _ := gzip.NewWriterLevel(buff, gzip.BestCompression)
	_, err = w.Write(data)
	if err != nil {
		b.Fatalf("unable to write: %v", err)
	}

	err = w.Close()
	if err != nil {
		b.Fatalf("unable to close: %v", err)
	}

	data = buff.Bytes()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r, _ := gzip.NewReader(bytes.NewReader(data))
		ungzipped, _ := ioutil.ReadAll(r)

		_, err := UnmarshalChangeset(ungzipped)
		if err != nil {
			b.Fatalf("unable to unmarshal: %v", err)
		}
	}
}

func loadChange(t testing.TB, filename string) *Change {
	data := readFile(t, filename)

	c := &Change{}
	err := xml.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("unable to unmarshal %s: %v", filename, err)
	}

	cleanXMLNameFromChange(c)
	return c
}

func loadOSM(t testing.TB, filename string) *OSM {
	data := readFile(t, filename)

	o := &OSM{}
	err := xml.Unmarshal(data, &o)
	if err != nil {
		t.Fatalf("unable to unmarshal %s: %v", filename, err)
	}

	cleanXMLNameFromOSM(o)
	return o
}

func readFile(t testing.TB, filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("unable to open %s: %v", filename, err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("unable to read file %s: %v", filename, err)
	}

	return data
}
