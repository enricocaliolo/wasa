export class ConversationParticipant {
  public name: string
  public userId: number
  public joinedAt: Date

  constructor(data: ConversationParticipantDTO) {
    this.name = data.name
    this.userId = data.user_id
    this.joinedAt = new Date(data.joined_at)
  }

  static fromJSON(json: ConversationParticipantDTO): ConversationParticipant {
    return new ConversationParticipant(json)
  }

  toJSON(): ConversationParticipantDTO {
    return {
      name: this.name,
      user_id: this.userId,
      joined_at: this.joinedAt.toISOString(),
    }
  }
}
export interface ConversationParticipantDTO {
  name: string
  user_id: number
  joined_at: string
}
