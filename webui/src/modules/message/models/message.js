import { User } from '../../auth/models/user';
import { Reaction } from './reaction';
import { RepliedToMessage } from './repliedToMessage';
export class Message {
	constructor(data) {
		this.messageId = data.message_id;
		this.content = data.content;
		this.contentType = data.content_type;
		this.sentTime = new Date(data.sent_time);
		this.editedTime = data.edited_time ? new Date(data.edited_time) : undefined;
		this.deletedTime = data.deleted_time
			? new Date(data.deleted_time)
			: undefined;
		this.conversationId = data.conversation_id;
		this.reactions = data.reactions
			? data.reactions.map((r) => new Reaction(r))
			: [];
		this.sender = new User(data.sender);
		this.repliedToMessage = Object.prototype.hasOwnProperty.call(
			data,
			'replied_to_message'
		)
			? new RepliedToMessage(data.replied_to_message)
			: null;
		this.isForwarded = data.is_forwarded;
		this.seenBy = Array.isArray(data.seen_by) ? data.seen_by : [];
	}

	isSeenBy(userId) {
		return Array.isArray(this.seenBy) && this.seenBy.includes(userId);
	}

	addSeenBy(userId) {
		if (!Array.isArray(this.seenBy)) {
			this.seenBy = [];
		}
		if (!this.seenBy.includes(userId)) {
			this.seenBy.push(userId);
		}
	}

	get seenCount() {
		return Array.isArray(this.seenBy) ? this.seenBy.length : 0;
	}

	static fromJSON(json) {
		return new Message(json);
	}

	toJSON() {
		return {
			message_id: this.messageId,
			content: this.content,
			content_type: this.contentType,
			sent_time: this.sentTime.toISOString(),
			edited_time: this.editedTime?.toISOString(),
			deleted_time: this.deletedTime?.toISOString(),
			conversation_id: this.conversationId,
			created_at: this.sentTime.toISOString(),
			sender: this.sender.toJSON(),
			replied_to_message: this.repliedToMessage?.toJSON(),
			is_forwarded: this.isForwarded,
			seen_by: this.seenBy,
		};
	}

	get isImage() {
		return this.contentType === 'image';
	}

	get displayContent() {
		if (this.isImage) {
			return `data:image/jpeg;base64,${this.content}`;
		}
		return this.content;
	}

	isEdited() {
		return !!this.editedTime;
	}

	isDeleted() {
		return !!this.deletedTime;
	}
}
