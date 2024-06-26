import { UpdateUserRequest } from "./UpdateUserRequest";
import type { User } from "./User";

 /**
 * @description OK
*/
export type PutApiUsers200 = User;

 /**
 * @description User data
*/
export type PutApiUsersMutationRequest = UpdateUserRequest;

 /**
 * @description OK
*/
export type PutApiUsersMutationResponse = User;

 export type PutApiUsersMutation = {
    Response: PutApiUsersMutationResponse;
    Request: PutApiUsersMutationRequest;
};