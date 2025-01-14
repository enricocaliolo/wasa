import api from "../../../shared/api/api";
import { Conversation } from "../models/conversation";
import { Message } from "../../message/models/message";

export const conversationAPI = {
	getUserConversations: async () => {
		const response = await api.get("/conversations");

		if (response.status === 200 && response.data != null) {
			return response.data.map((json) => Conversation.fromJSON(json));
		}
		return null
	},
	getConversation: async (conversation_id) => {
		const response = await api.get(`/conversations/${conversation_id}`);
		if (response.data == null) {
			return null;
		}
		var check = response.data.map((json) => new Message(json));
		return check;
	},
	createConversation: async (members, name) => {
		const response = await api.post("/conversations", {
			members: members,
			name: name,
		});
		return new Conversation(response.data);
	},
	updateGroupName: async (conversation_id, name) => {
		const response = await api.put(
			`/conversations/${conversation_id}/name`,
			{
				name: name,
			},
		);
		if (response.status === 200) {
			return true;
		} else {
			return false;
		}
	},
	updateGroupPhoto: async (conversation_id, photo) => {

		try {
			let payload = photo;
			let config = {
				headers: {
					"Content-Type": "image/*",
				},
			};
	
			// If the icon is a base64 string, convert it to a File object
			if (typeof photo === "string") {
				payload = imageConverter.base64ToFile(photo);
				if (!payload) {
					throw new Error("Failed to convert image");
				}
			}
	
			const response = await api.put(`/conversations/${conversation_id}/photo`, payload, config);
			if (response.status === 200) {
				return true
			}
			return false;
		} catch (error) {
			console.error("Error changing icon:", error);
			throw error;
		}
	},
	leaveGroup: async (conversation_id ) => {
		try{
			const response = await api.delete(`/conversations/${conversation_id}`);
		if (response.status === 200) {
			return true;
		}} catch(e) {
			if (e.response.status === 403) {
				throw new Error("Can't delete a direct conversation.")
			}
		}
	},
	addGroupMembers: async (conversationId, members) => {
		const response = await api.put(`/conversations/${conversationId}/users`, {
			members: members
		});
		
		if (response.status === 200) {
			return true;
		} else {
			throw new Error('Failed to add members to group');
		}
	}
};
