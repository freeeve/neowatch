package neowatch

import (
	"bytes"
	"fmt"
	"os"
	"syscall"
	"time"
   "math/rand"
)

var NIL_IDX uint64 = uint64(4294967295)
var NODE_SIZE = 14

type NodeRec struct {
	inUse        bool
	firstRel     uint64
	firstProp    uint64
	packedLabels uint64
}

func (node *NodeRec) Read(idx uint64, mmap []byte) {
	base := idx * 14
	node.inUse = (mmap[base] & 0x1) == 1
	node.firstRel = (uint64(mmap[base]&0xE) << 31) | uint64(mmap[base+1])<<24 | uint64(mmap[base+2])<<16 | uint64(mmap[base+3])<<8 | uint64(mmap[base+4])
	node.firstProp = (uint64(mmap[base]&0xF0) << 28) | uint64(mmap[base+5])<<24 | uint64(mmap[base+6])<<16 | uint64(mmap[base+7])<<8 | uint64(mmap[base+8])
	node.packedLabels = uint64(mmap[base+9])<<28 | uint64(mmap[base+10])<<24 | uint64(mmap[base+11])<<16 | uint64(mmap[base+12])<<8 | uint64(mmap[base+13])
}

func (node NodeRec) String() string {
	buf := fmt.Sprintf("inUse: %t\n", node.inUse)
	if node.firstRel == NIL_IDX {
		buf = buf + "firstRel: nil\n"
	} else {
		buf = buf + fmt.Sprintf("firstRel: %d\n", node.firstRel)
	}
	if node.firstProp == NIL_IDX {
		buf = buf + "firstProp: nil\n"
	} else {
		buf = buf + fmt.Sprintf("firstProp: %d\n", node.firstProp)
	}
	buf = buf + fmt.Sprintf("packedLabels: 0x%010X\n", node.packedLabels)
	return buf
}

func NodeStoreWatch(neopath string, ch chan string) {
	file, err := os.Open(neopath + "/neostore.nodestore.db")
	fd := int(file.Fd())
	mmap, err := syscall.Mmap(fd, 0, 100 * NODE_SIZE, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	node := NodeRec{}
	mmap_slice := make([]byte, len(mmap))
	copy(mmap_slice, mmap)
	go func() {
		for {
			for i := 0; i < 100; i++ {
				if !bytes.Equal(mmap[i*NODE_SIZE:i*NODE_SIZE+NODE_SIZE], mmap_slice[i*NODE_SIZE:i*NODE_SIZE+NODE_SIZE]) {
					node.Read(uint64(i), mmap)
					fmt.Printf("updated node idx %d: \n", i)
					fmt.Println(node)
				}
			}
			copy(mmap_slice, mmap)
			time.Sleep(time.Millisecond * time.Duration(rand.Int() % 100))
		}
	}()
}
