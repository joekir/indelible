package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"unsafe"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/joekir/indelible/rpc/indelible"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	_FS_IOC_SETFLAGS = uintptr(0x40086602) // assuming 64bit OS, TODO add startup check
	sockname         = "indelible.sock"
)

const (
	// from /usr/include/linux/fs.h
	FS_SECRM_FL        = 0x00000001 /* Secure deletion */
	FS_UNRM_FL         = 0x00000002 /* Undelete */
	FS_COMPR_FL        = 0x00000004 /* Compress file */
	FS_SYNC_FL         = 0x00000008 /* Synchronous updates */
	FS_IMMUTABLE_FL    = 0x00000010 /* Immutable file */
	FS_APPEND_FL       = 0x00000020 /* writes to file may only append */
	FS_NODUMP_FL       = 0x00000040 /* do not dump file */
	FS_NOATIME_FL      = 0x00000080 /* do not update atime */
	FS_DIRTY_FL        = 0x00000100
	FS_COMPRBLK_FL     = 0x00000200 /* One or more compressed clusters */
	FS_NOCOMP_FL       = 0x00000400 /* Don't compress */
	FS_ECOMPR_FL       = 0x00000800 /* Compression error */
	FS_BTREE_FL        = 0x00001000 /* btree format dir */
	FS_INDEX_FL        = 0x00001000 /* hash-indexed directory */
	FS_IMAGIC_FL       = 0x00002000 /* AFS directory */
	FS_JOURNAL_DATA_FL = 0x00004000 /* Reserved for ext3 */
	FS_NOTAIL_FL       = 0x00008000 /* file tail should not be merged */
	FS_DIRSYNC_FL      = 0x00010000 /* dirsync behaviour (directories only) */
	FS_TOPDIR_FL       = 0x00020000 /* Top of directory hierarchies*/
	FS_EXTENT_FL       = 0x00080000 /* Extents */
	FS_DIRECTIO_FL     = 0x00100000 /* Use direct i/o */
	FS_NOCOW_FL        = 0x00800000 /* Do not cow file */
	FS_PROJINHERIT_FL  = 0x20000000 /* Create with parents projid */
	FS_RESERVED_FL     = 0x80000000 /* reserved for ext2 lib */
)

type server struct{}

func (s *server) CreateLog(ctx context.Context, in *pb.LogFileRequest) (*empty.Empty, error) {
	if _, err := os.Stat(in.Path); os.IsNotExist(err) {
		return &empty.Empty{}, status.Errorf(codes.InvalidArgument, "Log file does not exist, or indelible cannot read it")
	}

	log.Printf("Attempting to change file: %s to appendOnly", in.Path)
	err := setAppendOnly(in.Path)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.InvalidArgument, "unable to change attributes")
	}

	log.Printf("Successfully set file: %s to appendOnly", in.Path)
	return &empty.Empty{}, nil
}

func cleanup(socketpath string) {
	log.Printf("Unlinking sock file: %s", socketpath)
	syscall.Unlink(socketpath)
}

func main() {
	socketpath := path.Join("/tmp", sockname)

	// Remove previously existing sock file
	syscall.Unlink(socketpath)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(socketpath)
		os.Exit(1)
	}()

	log.Printf("Linking sock file: %s", socketpath)
	lis, err := net.Listen("unix", socketpath)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLogCreatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Leveraged examples from https://github.com/snapcore/snapd/blob/master/osutil/chattr.go
func setAppendOnly(path string) error {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return err
	}

	// cannot use the const, need a stack varaiable ptr here
	var attr int32 = int32(FS_APPEND_FL | FS_EXTENT_FL)
	return ioctl(f, _FS_IOC_SETFLAGS, &attr)
}

func ioctl(f *os.File, request uintptr, attrp *int32) error {
	argp := uintptr(unsafe.Pointer(attrp))
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), request, argp)
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}

	return nil
}
