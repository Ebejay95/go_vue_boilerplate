import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { UserServiceClient } from "./generated/user.client";

const transport = new GrpcWebFetchTransport({
	baseUrl: "http://localhost:8080",
});

export const userServiceClient = new UserServiceClient(transport);

// Fallback: Wenn gRPC-Web nicht verfügbar, nutze REST
export class UserGrpcService {
	async getUsers() {
		try {
			// Hier würde normalerweise gRPC-Web Code stehen
			// Für Einfachheit nutzen wir erstmal REST als Fallback
			const response = await fetch('http://localhost:8080/api/users');
			const users = await response.json();
			return { users };
		} catch (error) {
			console.error('gRPC call failed, using REST fallback:', error);
			throw error;
		}
	}

	async createUser(name, email) {
		try {
			const response = await fetch('http://localhost:8080/api/users', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({ name, email }),
			});
			return await response.json();
		} catch (error) {
			console.error('gRPC call failed:', error);
			throw error;
		}
	}

	async deleteUser(id) {
		try {
			const response = await fetch(`http://localhost:8080/api/users/${id}`, {
				method: 'DELETE',
			});
			return await response.json();
		} catch (error) {
			console.error('gRPC call failed:', error);
			throw error;
		}
	}
}