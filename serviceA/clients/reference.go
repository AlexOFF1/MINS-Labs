package clients

import (
	"context"
	"fmt"
	"time"

	pb "mins_EduCenter/proto/reference"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ReferenceClient struct {
	conn   *grpc.ClientConn
	client pb.ReferenceServiceClient
}

func NewReferenceClient(addr string) (*ReferenceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &ReferenceClient{
		conn:   conn,
		client: pb.NewReferenceServiceClient(conn),
	}, nil
}

func (c *ReferenceClient) Close() error {
	return c.conn.Close()
}

func (c *ReferenceClient) GetCourse(ctx context.Context, courseID string) (*pb.Course, error) {
	// Извлекаем traceID из контекста
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "unknown"
	}
	md := metadata.New(map[string]string{"trace-id": traceID.(string)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := c.client.GetCourse(ctx, &pb.GetCourseRequest{CourseId: courseID})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("course %s not found in reference service", courseID)
			case codes.Unavailable, codes.DeadlineExceeded:
				return nil, fmt.Errorf("reference service is unavailable, please try again later")
			default:
				return nil, fmt.Errorf("reference service error: %v", err)
			}
		}
		return nil, fmt.Errorf("failed to call reference service: %v", err)
	}
	return resp, nil
}
