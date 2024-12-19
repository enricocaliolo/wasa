import type { User } from '@/modules/auth/models/user'
import type { Conversation } from '@/modules/conversation/models/conversation'
import type { Reaction } from './reaction'

export class Message {
  public messageId: number
  public content: Blob
  public contentType: string
  public sentTime: Date
  public editedTime?: Date
  public deletedTime?: Date
  public senderId: number
  public conversationId: number
  public repliedTo?: number
  public forwardedFrom?: number
  public reactions?: Reaction[]
  public sender: User

  constructor(data: MessageDTO) {
    this.messageId = data.message_id
    this.content = data.content
    this.contentType = data.content_type
    this.sentTime = new Date(data.sent_time)
    this.editedTime = data.edited_time ? new Date(data.edited_time) : undefined
    this.deletedTime = data.deleted_time ? new Date(data.deleted_time) : undefined
    this.senderId = data.sender_id
    this.conversationId = data.conversation_id
    this.repliedTo = data.replied_to
    this.forwardedFrom = data.forwarded_from
    this.sender = data.sender
  }

  static fromJSON(json: MessageDTO): Message {
    return new Message(json)
  }

  toJSON(): MessageDTO {
    return {
      message_id: this.messageId,
      content: this.content,
      content_type: this.contentType,
      sent_time: this.sentTime.toISOString(),
      edited_time: this.editedTime?.toISOString(),
      deleted_time: this.deletedTime?.toISOString(),
      sender_id: this.senderId,
      conversation_id: this.conversationId,
      replied_to: this.repliedTo,
      forwarded_from: this.forwardedFrom,
      created_at: this.sentTime.toISOString(),
      sender: this.sender,
    }
  }

  isEdited(): boolean {
    return !!this.editedTime
  }

  isDeleted(): boolean {
    return !!this.deletedTime
  }
}

export interface MessageDTO {
  message_id: number
  content: Blob
  content_type: string
  sent_time: string
  edited_time?: string
  deleted_time?: string
  sender_id: number
  conversation_id: number
  replied_to?: number
  forwarded_from?: number
  created_at: string
  sender: User
}
