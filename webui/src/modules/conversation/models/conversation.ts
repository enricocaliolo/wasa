import { Message, type MessageDTO } from '@/modules/message/models/message'
import { ConversationParticipant } from './conversation_participant'

export class Conversation {
  public conversationId: number
  public name: string
  public isGroup: boolean
  public createdAt: Date
  public participants: ConversationParticipant[]
  public messages?: Message[]

  constructor(data: ConversationDTO) {
    this.conversationId = data.conversation_id
    this.name = data.name
    this.isGroup = data.is_group
    this.createdAt = new Date(data.created_at)
    this.participants = data.conversation_participants
    this.messages = data.messages ? data.messages.map((messageDto) => new Message(messageDto)) : []
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
      conversation_participants: this.participants,
      messages: this.messages?.map((message) => message.toJSON()) || [],
    }
  }
}

export interface ConversationDTO {
  conversation_id: number
  name: string
  is_group: boolean
  created_at: string
  conversation_participants: ConversationParticipant[]
  messages?: MessageDTO[]
}
