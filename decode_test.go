package osmpbf

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	// originally downloaded from http://download.geofabrik.de/europe/great-britain/england/greater-london.html
	London    = "greater-london-140324.osm.pbf"
	LondonURL = "https://googledrive.com/host/0B8pisLiGtmqDR3dOR3hrWUpRTVE"
)

func init() {
	_, err := os.Stat(London)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("\nDownload %s from %s.\nFor example: 'wget -O %s %s'", London, LondonURL, London, LondonURL))
	}
}

func TestDecoder(t *testing.T) {
	f, err := os.Open(London)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var n *Node
	en := &Node{
		ID:  18088578,
		Lat: 51.5442632,
		Lon: -0.2010027,
		Tags: map[string]string{
			"alt_name":   "The King's Head",
			"amenity":    "pub",
			"created_by": "JOSM",
			"name":       "The Luminaire",
			"note":       "Live music venue too",
		},
	}

	var w *Way
	ew := &Way{
		ID:      4257116,
		NodeIDs: []int64{21544864, 333731851, 333731852, 333731850, 333731855, 333731858, 333731854, 108047, 769984352, 21544864},
		Tags: map[string]string{
			"area":    "yes",
			"highway": "pedestrian",
			"name":    "Fitzroy Square",
		},
	}

	var r *Relation
	er := &Relation{
		ID: 7677,
		Members: []Member{
			Member{ID: 4875932, Type: WayType, Role: "outer"},
			Member{ID: 4894305, Type: WayType, Role: "inner"},
		},
		Tags: map[string]string{
			"created_by": "Potlatch 0.9c",
			"type":       "multipolygon",
		},
	}

	idsExpectedOrder := []string{
		"node/44", "node/47", "node/52", "node/58", "node/60", // start of dense nodes
		"node/79",                                                                                     // just because way/79 is already there
		"node/2740703694", "node/2740703695", "node/2740703697", "node/2740703699", "node/2740703701", // end of dense nodes
		"way/73", "way/74", "way/75", "way/79", "way/482", // start of ways
		"way/268745428", "way/268745431", "way/268745434", "way/268745436", "way/268745439", // end of ways
		"relation/69", "relation/94", "relation/152", "relation/245", "relation/332", // start of relations
		"relation/3593436", "relation/3595575", "relation/3595798", "relation/3599126", "relation/3599127", // end of relations
	}
	ids := make(map[string]bool)
	for _, id := range idsExpectedOrder {
		ids[id] = true
	}

	var nc, wc, rc int
	var id string
	enc, ewc, erc := 2729006, 459055, 12833
	idsOrder := make([]string, 0, len(idsExpectedOrder))
	d := NewDecoder(f)
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		} else {
			switch v := v.(type) {
			case *Node:
				nc++
				if v.ID == en.ID {
					n = v
				}
				id = fmt.Sprintf("node/%d", v.ID)
				if ids[id] {
					idsOrder = append(idsOrder, id)
				}
			case *Way:
				wc++
				if v.ID == ew.ID {
					w = v
				}
				id = fmt.Sprintf("way/%d", v.ID)
				if ids[id] {
					idsOrder = append(idsOrder, id)
				}
			case *Relation:
				rc++
				if v.ID == er.ID {
					r = v
				}
				id = fmt.Sprintf("relation/%d", v.ID)
				if ids[id] {
					idsOrder = append(idsOrder, id)
				}
			default:
				t.Fatalf("unknown type %T", v)
			}
		}
	}

	if !reflect.DeepEqual(en, n) {
		t.Errorf("\nExpected: %#v\nActual:   %#v", en, n)
	}
	if !reflect.DeepEqual(ew, w) {
		t.Errorf("\nExpected: %#v\nActual:   %#v", ew, w)
	}
	if !reflect.DeepEqual(er, r) {
		t.Errorf("\nExpected: %#v\nActual:   %#v", er, r)
	}
	if enc != nc || ewc != wc || erc != rc {
		t.Errorf("\nExpected %7d nodes, %7d ways, %7d relations\nGot      %7d nodes, %7d ways, %7d relations", enc, ewc, erc, nc, wc, rc)
	}
	if !reflect.DeepEqual(idsExpectedOrder, idsOrder) {
		t.Errorf("\nExpected: %v\nGot:      %v", idsExpectedOrder, idsOrder)
	}
}

func BenchmarkDecoder(b *testing.B) {
	file := os.Getenv("OSMPBF_BENCHMARK_FILE")
	if file == "" {
		file = London
	}
	f, err := os.Open(file)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		d := NewDecoder(f)
		n, w, r, c, start := 0, 0, 0, 0, time.Now()
		for {
			if v, err := d.Decode(); err == io.EOF {
				break
			} else if err != nil {
				b.Fatal(err)
			} else {
				switch v := v.(type) {
				case *Node:
					n++
				case *Way:
					w++
				case *Relation:
					r++
				default:
					b.Fatalf("unknown type %T", v)
				}
			}
			c++
		}

		b.Logf("Done in %.3f seconds. Total: %d, Nodes: %d, Ways: %d, Relations: %d\n",
			time.Now().Sub(start).Seconds(), c, n, w, r)
	}
}
