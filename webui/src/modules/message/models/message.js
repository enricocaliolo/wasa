import { User } from '@/modules/auth/models/user'

export class Message {
  constructor(data) {
    this.messageId = data.message_id
    this.content = data.content
    this.contentType = data.content_type
    this.sentTime = new Date(data.sent_time)
    this.editedTime = data.edited_time ? new Date(data.edited_time) : undefined
    this.deletedTime = data.deleted_time ? new Date(data.deleted_time) : undefined
    this.conversationId = data.conversation_id
    this.repliedTo = data.replied_to
    this.forwardedFrom = data.forwarded_from
    this.reactions = data.reactions
    this.sender = new User(data.sender)
  }

  static fromJSON(json) {
    return new Message(json)
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
      replied_to: this.repliedTo,
      forwarded_from: this.forwardedFrom,
      created_at: this.sentTime.toISOString(),
      sender: this.sender.toJSON(),
    }
  }

  isEdited() {
    return !!this.editedTime
  }

  isDeleted() {
    return !!this.deletedTime
  }
}