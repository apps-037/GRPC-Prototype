import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import org.springframework.boot.CommandLineRunner;

import java.util.Scanner;

import javax.annotation.PreDestroy;
import com.example.grpc.Test;
import com.example.grpc.Test.*;

@SpringBootApplication
public class GrpcClientApplication {

    public static void main(String[] args) {
        SpringApplication.run(GrpcClientApplication.class, args);
    }

    @Bean
    public ManagedChannel managedChannel() {
        return ManagedChannelBuilder.forAddress("localhost", 50051)
                .usePlaintext()
                .build();
    }

    @PreDestroy
    public void shutdownChannel(ManagedChannel channel) {
        if (channel != null && !channel.isShutdown()) {
            channel.shutdown();
        }
    }

    @Component
    public static class GrpcClientRunner implements CommandLineRunner {

        private final Test myServiceStub;
        private String name;
        private int age;

        public GrpcClientRunner(ManagedChannel managedChannel) {
            myServiceStub = Test.newBlockingStub(managedChannel); //the method calls made on the stub will block until a response is received.
        }

        @Override
        public void run(String... args) {

        Scanner scanner = new Scanner(System.in);

        System.out.println("Enter name: ");
         name = scanner.nextLine();
    
        System.out.print("Enter age: ");
         age = scanner.nextInt();
            // Calling gRPC methods here
            createRecord();
            updateRecord();
        }

        private void createRecord() {    //CreateRecordRequest obj is created to sent data to grpc server 
            CreateRecordRequest request = ((CreateRecordRequest.Builder) CreateRecordRequest.newBuilder())
                    .setName(name)
                    .setAge(age)
                    .build();

            CreateRecordResponse response = myServiceStub.createRecord(request);
            System.out.println("Created record with ID: " + response.getId());
        }

        private void updateRecord() {
            UpdateRecordRequest request = UpdateRecordRequest.newBuilder()
                    .setId("insert_record_id_here")
                    .setName("Updated Name")
                    .setAge(35)
                    .build();

            UpdateRecordResponse response = myServiceStub.updateRecord(request);
            if (response.getSuccess()) {
                System.out.println("Record updated successfully.");
            } else {
                System.out.println("Failed to update record.");
            }
        }
    }
}


//In summary, the GrpcClientRunner component is responsible for invoking gRPC methods (createRecord and updateRecord) using the 
// myServiceStub object, which is created with the ManagedChannel dependency. It demonstrates how to integrate gRPC client functionality 
// into a Spring Boot application and execute gRPC operations during application startup.

// a stub is being used to invoke the gRPC methods on the server. A stub is a client-side representation of the gRPC service that 
// allows the client to make remote procedure calls (RPCs) to the server.

// In gRPC, the client and server communicate using protocol buffers (protobuf) messages and define the service interface in a .proto 
// file. The gRPC code generator generates client-side stubs based on the service definition, which encapsulates the details of sending 
// requests and receiving responses over the network.


