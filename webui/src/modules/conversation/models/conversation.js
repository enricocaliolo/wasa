// import { Message } from '@/modules/message/models/message'

export class Conversation {
	constructor(data) {
		this.conversationId = data.conversation_id;
		this.name = data.name;
		this.isGroup = data.is_group;
		this.createdAt = new Date(data.created_at);
		// this.participants = data.participants.map(participant => participant.user_id)
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
			is_group: this.isGroup,
			created_at: this.createdAt.toISOString(),
			// conversation_participants: this.participants,
			messages: this.messages?.map((message) => message.toJSON()) || [],
		};
	}
}
