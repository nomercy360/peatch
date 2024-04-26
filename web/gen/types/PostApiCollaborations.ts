import { CreateCollaboration } from './CreateCollaboration';
import type { Collaboration } from './Collaboration';

/**
 * @description Created
*/
export type PostApiCollaborations201 = Collaboration;

 /**
 * @description Collaboration data
*/
 export type PostApiCollaborationsMutationRequest = CreateCollaboration;

 /**
 * @description Created
*/
export type PostApiCollaborationsMutationResponse = Collaboration;

 export type PostApiCollaborationsMutation = {
    Response: PostApiCollaborationsMutationResponse;
    Request: PostApiCollaborationsMutationRequest;
};