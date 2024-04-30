export type CreateUserCollaboration = {
    /**
     * @type string | undefined
    */
    message?: string;
    /**
     * @type integer
    */
    requester_id: number;
    /**
     * @type integer
    */
    user_id: number;
};