import { User } from "./User";

 export type UserWithToken = {
    /**
     * @type string | undefined
    */
    token?: string;
    user?: User;
};