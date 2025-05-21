import { grpc } from "@improbable-eng/grpc-web";
import { BrowserHeaders } from "browser-headers";

// Import the interfaces from our generated code
import { 
  PostRequest, 
  CommentsResponse,
  Comment 
} from "../generated/reddit_service";

/**
 * Type definition for the Reddit service interface
 */
export interface RedditService {
  GetPostComments(request: PostRequest): Promise<CommentsResponse>;
}

/**
 * Error class for gRPC-specific errors
 */
export class GrpcError extends Error {
  code: grpc.Code;
  
  constructor(message: string, code: grpc.Code) {
    super(message);
    this.name = "GrpcError";
    this.code = code;
  }
}

// Add serializeBinary directly to request/response objects to match expectations
// This is a workaround since we're missing the expected methods from gRPC-web
(PostRequest as any).serializeBinary = function(instance: PostRequest) {
  // Simple binary serialization - just JSON for now
  const json = JSON.stringify(instance);
  console.log("Serializing request:", json);
  return new TextEncoder().encode(json);
};

(CommentsResponse as any).deserializeBinary = function(bytes: Uint8Array) {
  try {
    // Create response from JSON
    const json = new TextDecoder().decode(bytes);
    console.log("Deserializing response:", json);
    const data = JSON.parse(json);
    
    return {
      postTitle: data.postTitle || "Received response",
      comments: Array.isArray(data.comments) 
        ? data.comments.map((c: any) => ({ body: c.body || "" }))
        : []
    };
  } catch (error) {
    console.error("Error deserializing response:", error);
    return {
      postTitle: "Error parsing response",
      comments: [{ body: "Failed to parse response" }]
    };
  }
};

// Service definition
const SERVICE = {
  serviceName: "reddit_implementation.RedditService"
};

/**
 * Method definition that matches exactly what grpc-web expects
 */
const methodDescriptor: grpc.MethodDefinition<PostRequest, CommentsResponse> = {
  methodName: "GetPostComments",
  service: SERVICE,
  requestStream: false,
  responseStream: false,
  requestType: PostRequest as any,
  responseType: CommentsResponse as any
};

/**
 * Creates and returns a gRPC-web client for the RedditService
 * 
 * @param host The gRPC server host URL (e.g., http://localhost:8080)
 * @returns A RedditService client with typed methods
 */
export const createRedditServiceClient = (host: string): RedditService => {
  // Validate the host URL
  if (!host) {
    throw new Error("gRPC host URL is required");
  }
  
  if (!host.startsWith("http://") && !host.startsWith("https://")) {
    throw new Error("gRPC host URL must start with http:// or https://");
  }
  
  return {
    /**
     * Get Reddit post comments based on search query
     * 
     * @param request The PostRequest message with subreddit, search query, and limit
     * @returns Promise resolving to CommentsResponse
     */
    GetPostComments: (request: PostRequest): Promise<CommentsResponse> => {
      return new Promise((resolve, reject) => {
        console.log("Making gRPC-web request to:", host);
        console.log("Request data:", request);
        
        try {
          // Execute the unary gRPC call
          grpc.unary(methodDescriptor, {
            request,
            host,
            // Set proper headers for gRPC-web
            metadata: new BrowserHeaders({
              "Content-Type": "application/grpc-web+json", // Use JSON format instead of proto
              "Accept": "application/grpc-web+json"
            }),
            onEnd: (response) => {
              const { status, statusMessage, headers, message, trailers } = response;
              
              // Handle successful response
              if (status === grpc.Code.OK && message) {
                console.log("Received successful gRPC response");
                resolve(message as CommentsResponse);
              } 
              // Handle error cases
              else {
                const errorMessage = statusMessage || "Unknown gRPC error";
                console.error(`gRPC error: ${errorMessage} (code: ${status})`, {
                  status,
                  headers: headers?.toHeaders(),
                  trailers: trailers?.toHeaders()
                });
                
                reject(new GrpcError(errorMessage, status));
              }
            }
          });
        } catch (error) {
          console.error("Error making gRPC-web request:", error);
          reject(new Error(`Failed to make gRPC-web request: ${error}`));
        }
      });
    }
  };
};

/**
 * Utility function to detect if an error is a GrpcError
 */
export const isGrpcError = (error: unknown): error is GrpcError => {
  return error instanceof GrpcError;
};

