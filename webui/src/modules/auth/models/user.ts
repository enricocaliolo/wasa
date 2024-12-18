export class User {
  public userId: number
  public username: string
  public icon?: string
  public createdAt: Date

  constructor(data: UserDTO) {
    this.userId = data.user_id
    this.username = data.username
    this.icon = data.icon
    this.createdAt = new Date(data.created_at)
  }

  static fromJSON(json: UserDTO): User {
    return new User(json)
  }

  toJSON(): UserDTO {
    return {
      user_id: this.userId,
      username: this.username,
      icon: this.icon,
      created_at: this.createdAt.toISOString(),
    }
  }
}

export interface UserDTO {
  user_id: number
  username: string
  icon?: string
  created_at: string
}
