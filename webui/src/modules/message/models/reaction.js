import { User } from "@/modules/auth/models/user";

export class Reaction {
	constructor(data) {
		this.reactionId = data.reaction_id;
		this.messageId = data.message_id;
		this.reaction = data.reaction;
		this.user = new User(data.user);
	}

	static fromJSON(json) {
		return new Reaction(json);
	}

	toJSON() {
		return {
			reaction_id: this.reactionId,
			message_id: this.messageId,
			reaction: this.reaction,
			user: this.user.toJSON(),
		};
	}

	isEdited() {
		return !!this.editedTime;
	}

	isDeleted() {
		return !!this.deletedTime;
	}
}
