export class RepliedToMessage {
	constructor(data) {
		this.messageId = data.message_id;
		this.content = data.content;
		this.contentType = data.content_type;
	}
	static fromJSON(json) {
		return new RepliedToMessage(json);
	}

	toJSON() {
		return {
			message_id: this.messageId,
			content: this.content,
			content_type: this.contentType,
		};
	}
}
