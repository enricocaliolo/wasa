import type { User } from '@/modules/auth/models/user'
import type { Message } from './message'

export class Reaction {
  public reactionId: number
  public messageId: number
  public userId: number
  public reaction: Blob
  public user?: User
  public message?: Message

  constructor(data: ReactionDTO) {
    this.reactionId = data.reaction_id
    this.messageId = data.message_id
    this.userId = data.user_id
    this.reaction = data.reaction
  }

  static fromJSON(json: ReactionDTO): Reaction {
    return new Reaction(json)
  }

  toJSON(): ReactionDTO {
    return {
      reaction_id: this.reactionId,
      message_id: this.messageId,
      user_id: this.userId,
      reaction: this.reaction,
    }
  }
}

export interface ReactionDTO {
  reaction_id: number
  message_id: number
  user_id: number
  reaction: Blob
}
