import type { Message } from '@/modules/message/models/message'
import type { ConversationParticipant } from './conversation_participant'

export class Conversation {
  public conversationId: number
  public name?: string
  public isGroup: boolean
  public createdAt: Date
  public participants?: ConversationParticipant[]
  public messages?: Message[]

  constructor(data: ConversationDTO) {
    this.conversationId = data.conversation_id
    this.name = data.name
    this.isGroup = data.is_group
    this.createdAt = new Date(data.created_at)
  }

  static fromJSON(json: ConversationDTO): Conversation {
    return new Conversation(json)
  }

  toJSON(): ConversationDTO {
    return {
      conversation_id: this.conversationId,
      name: this.name,
      is_group: this.isGroup,
      created_at: this.createdAt.toISOString(),
    }
  }
}

export interface ConversationDTO {
  conversation_id: number
  name?: string
  is_group: boolean
  created_at: string
}
