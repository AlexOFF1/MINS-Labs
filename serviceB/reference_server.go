package main

import (
	"context"
	"log"

	pb "mins_EduCenter/proto/reference"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ReferenceServer struct {
	pb.UnimplementedReferenceServiceServer
	courses map[string]*pb.Course
}

func NewReferenceServer() *ReferenceServer {
	// Тестовые курсы
	courses := map[string]*pb.Course{
		"course-1": {Id: "course-1", Name: "Go_Basics", Description: "Learn Go", DurationWeeks: 4, Price: 20000},
		"course-2": {Id: "course-2", Name: "Python_Intro", Description: "Python basics", DurationWeeks: 3, Price: 15000},
		"course-3": {Id: "course-3", Name: "Java_Advanced", Description: "Advanced Java", DurationWeeks: 6, Price: 30000},
	}
	return &ReferenceServer{courses: courses}
}

func (s *ReferenceServer) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.Course, error) {
	course, ok := s.courses[req.CourseId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "course %s not found", req.CourseId)
	}
	return course, nil
}

func (s *ReferenceServer) ListCourses(ctx context.Context, req *pb.ListCoursesRequest) (*pb.ListCoursesResponse, error) {
	traceID := extractTraceID(ctx)
	log.Printf("[TraceID=%s] ListCourses called", traceID)

	courses := make([]*pb.Course, 0, len(s.courses))
	for _, c := range s.courses {
		courses = append(courses, c)
	}
	return &pb.ListCoursesResponse{Courses: courses}, nil
}

// Извлекаем trace-id из метаданных gRPC
func extractTraceID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "unknown"
	}
	if val := md["trace-id"]; len(val) > 0 {
		return val[0]
	}
	return "unknown"
}

// Интерсептор для логирования каждого запроса с trace-id
func TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	traceID := extractTraceID(ctx)
	log.Printf("[TraceID=%s] %s", traceID, info.FullMethod)
	return handler(ctx, req)
}
