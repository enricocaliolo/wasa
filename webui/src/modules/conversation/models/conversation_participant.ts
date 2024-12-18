export class ConversationParticipant {
  public conversationId?: number
  public userId: number
  public joinedAt: Date

  constructor(data: ConversationParticipantDTO) {
    this.conversationId = data.conversation_id
    this.userId = data.user_id
    this.joinedAt = new Date(data.joined_at)
  }

  static fromJSON(json: ConversationParticipantDTO): ConversationParticipant {
    return new ConversationParticipant(json)
  }

  toJSON(): ConversationParticipantDTO {
    return {
      conversation_id: this.conversationId,
      user_id: this.userId,
      joined_at: this.joinedAt.toISOString(),
    }
  }
}
export interface ConversationParticipantDTO {
  conversation_id?: number
  user_id: number
  joined_at: string
}
