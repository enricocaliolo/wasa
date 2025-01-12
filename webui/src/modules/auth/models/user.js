export class User {
	constructor(data) {
		this.userId = data.user_id;
		this.username = data.username;
		this.icon = data.icon;
		this.createdAt = data.created_at
			? new Date(data.created_at)
			: undefined;
	}

	static fromJSON(json) {
		return new User(json);
	}

	toJSON() {
		return {
			user_id: this.userId,
			username: this.username,
			icon: this.icon,
			created_at: this.createdAt?.toISOString(),
		};
	}

	get displayIcon() {
		if (this.icon) {
			return `data:image/jpeg;base64,${this.icon}`;
		}
	}
}
