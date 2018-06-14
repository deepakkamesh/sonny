package devices

import (
	"time"

	"github.com/blackjack/webcam"
	"github.com/golang/glog"
	"github.com/saljam/mjpeg"
)

// V4L format identifiers from /usr/include/linux/videodev2.h.
const (
	MJPEG   webcam.PixelFormat = 1196444237
	YUYV422 webcam.PixelFormat = 1448695129
)

// Width, Height.
var CamResolutions = map[int][]int{
	1:  {160, 120},
	2:  {176, 144},
	3:  {320, 176},
	4:  {320, 240},
	5:  {352, 288},
	6:  {432, 240},
	7:  {544, 288},
	8:  {640, 360},
	9:  {640, 480},
	10: {800, 480},
	11: {1024, 768},
}

type Video struct {
	Stream      *mjpeg.Stream
	cam         *webcam.Webcam
	height      uint32
	width       uint32
	pixelFormat webcam.PixelFormat
	stop        chan struct{}
	fps         uint
	capStatus   bool
}

func NewVideo(pixelFormat webcam.PixelFormat, w uint32, h uint32, fps uint) *Video {
	return &Video{
		pixelFormat: pixelFormat,
		height:      h,
		width:       w,
		stop:        make(chan struct{}),
		fps:         fps,
		capStatus:   false,
		Stream:      mjpeg.NewStream(),
	}
}

func (s *Video) GetVideoStream() *mjpeg.Stream {
	return s.Stream
}

func (s *Video) SetResMode(i int) {
	s.SetRes(uint32(CamResolutions[i][0]), uint32(CamResolutions[i][1]))
}

func (s *Video) SetRes(w uint32, h uint32) {
	s.height = h
	s.width = w
}

func (s *Video) SetFPS(fps uint) {
	s.fps = fps
}

func (s *Video) StartVideoStream() error {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		return err
	}

	if _, _, _, err := cam.SetImageFormat(s.pixelFormat, s.width, s.height); err != nil {
		return err
	}

	s.cam = cam

	if !s.capStatus {
		go s.startStreamer()
		return nil
	}
	glog.Info("Video capture already running")
	return nil
}

func (s *Video) StopVideoStream() {
	if s.capStatus {
		s.stop <- struct{}{}
		if err := s.cam.Close(); err != nil {
			glog.Errorf("Failed to stop stream:%v", err)
		}
	}
}

func (s *Video) startStreamer() {

	// Since the ReadFrame is buffered, trying to read at FPS results in delay.
	fpsTicker := time.NewTicker(time.Duration(1000/s.fps) * time.Millisecond)

	if err := s.cam.StartStreaming(); err != nil {
		glog.Errorf("Failed to start stream:%v", err)
		return
	}
	s.capStatus = true
	glog.Infof("Started Video Capture")

	frame := []byte{}
	for {
		select {
		case <-s.stop:
			glog.Info("Stopped Video Capture")
			s.capStatus = false
			return

		default:
			if err := s.cam.WaitForFrame(5); err != nil {
				glog.Errorf("Failed to read webcam:%v", err)
			}
			var err error
			frame, err = s.cam.ReadFrame()
			if err != nil || len(frame) == 0 {
				glog.Errorf("Failed tp read webcam frame:%v or frame size 0", err)
			}

		case <-fpsTicker.C:
			s.Stream.UpdateJPEG(frame)
		}
	}
}
