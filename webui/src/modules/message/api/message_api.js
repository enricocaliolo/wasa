import api from "../../../shared/api/api";
import { imageConverter } from "../helper/image_converter";

export const messagesAPI = {
	commentMessage: async (conversation_id, message_id, comment) => {
		const response = await api.put(
			`/conversations/${conversation_id}/messages/${message_id}`,
			{
				reaction: comment,
			},
		);
		if (response.status === 202) {
			return response.data;
		}
	},
	uncommentMessage: async (conversation_id, message_id, reaction_id) => {
		const response = await api.delete(
			`/conversations/${conversation_id}/messages/${message_id}/reactions/${reaction_id}`,
		);
		if (response.status === 202) {
			return true;
		}
	},
	sendMessage: async ({
		conversation_id,
		content,
		content_type = "text",
		replied_to_message = null,
		destination_conversation_id = null,
	}) => {
		try {
			let endpoint = `/conversations/${conversation_id}`;
			let payload = {};
			let config = {};

			if (content_type === "image") {
				let imageFile = content;
				if (typeof content === "string") {
					// I need to convert the image from base64 to a File object
					// because, otherwise, it was not working.
					imageFile = imageConverter.base64ToFile(content);
					if (!imageFile) {
						throw new Error("Failed to convert image");
					}
				}

				if (destination_conversation_id) {
					endpoint += `/forward?destination_conversation_id=${destination_conversation_id}`;
				} else if (replied_to_message) {
					endpoint += "/reply";
				}

				payload = imageFile;
				config = {
					headers: {
						"Content-Type": "image/*",
					},
				};
			} else {
				if (destination_conversation_id) {
					endpoint += "/forward";
					payload = {
						content,
						content_type,
						destination_conversation_id,
					};
				} else if (replied_to_message) {
					endpoint += "/reply";
					payload = {
						content,
						content_type,
						replied_to: replied_to_message,
					};
				} else {
					payload = {
						content,
						content_type,
					};
				}
			}

			const response = await api.post(endpoint, payload, config);

			if (response.status === 201) {
				return response.data;
			}
		} catch (error) {
			console.error("Error sending message:", error);
			throw error;
		}
	},
};
