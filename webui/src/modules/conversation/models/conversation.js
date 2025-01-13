// import { Message } from '@/modules/message/models/message'

import { User } from "../../auth/models/user";
import { Message } from "../../message/models/message";

export class Conversation {
	constructor(data) {
		this.conversationId = data.conversation_id;
		this.name = data.name;
		this.photo = data.photo;
		this.isGroup = data.is_group;
		this.createdAt = new Date(data.created_at);
		this.participants = data.participants.map((participant) => new User(participant))
		this.messages = data.messages
			? data.messages.map((messageDto) => new Message(messageDto))
			: [];
	}

	static fromJSON(json) {
		return new Conversation(json);
	}

	toJSON() {
		return {
			conversation_id: this.conversationId,
			name: this.name,
			photo: this.photo,
			is_group: this.isGroup,
			created_at: this.createdAt.toISOString(),
			conversation_participants: this.participants.map((participant) => participant.toJSON()),
			messages: this.messages.map((message) => message.toJSON()) || [],
		};
	}

	get displayPhoto() {
		if (this.photo) {
			return `data:image/jpeg;base64,${this.photo}`;
		}
		return 'user icon';
	}
}
