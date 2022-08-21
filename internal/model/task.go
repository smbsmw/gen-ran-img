package model

import "youtube-downloader-rest/pb"

type Task struct {
	Id  string `json:"id"`
	ImageName string `json:"image_name"`
}

func (t *Task) ToProto() *pb.Task {
	return &pb.Task{ImageName: t.ImageName}
}

func TaskFromProto(task *pb.Task) (*Task, error) {
	return &Task{
		ImageName: task.ImageName,
	}, nil
}
