// package main

// import "google.golang.org/grpc/test/bufconn"

// const bufSize = 1024 * 1024

// var lis *bufconn.Listener

// func init() {
//     lis = bufconn.Listen(bufSize)
//     s := grpc.NewServer()
//     pb.RegisterGreeterServer(s, &server{})
//     go func() {
//         if err := s.Serve(lis); err != nil {
//             log.Fatalf("Server exited with error: %v", err)
//         }
//     }()
// }

// func bufDialer(context.Context, string) (net.Conn, error) {
//     return lis.Dial()
// }

// func TestSayHello(t *testing.T) {
//     ctx := context.Background()
//     conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
//     if err != nil {
//         t.Fatalf("Failed to dial bufnet: %v", err)
//     }
//     defer conn.Close()
//     client := pb.NewGreeterClient(conn)
//     resp, err := client.SayHello(ctx, &pb.HelloRequest{"Dr. Seuss"})
//     if err != nil {
//         t.Fatalf("SayHello failed: %v", err)
//     }
//     log.Printf("Response: %+v", resp)
//     // Test for output here.