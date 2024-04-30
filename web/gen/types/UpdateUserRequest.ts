export type UpdateUserRequest = {
    /**
     * @type string
    */
    avatar_url: string;
    /**
     * @type array
    */
    badge_ids: number[];
    /**
     * @type string | undefined
    */
    city?: string;
    /**
     * @type string
    */
    country: string;
    /**
     * @type string
    */
    country_code: string;
    /**
     * @type string
    */
    description: string;
    /**
     * @type string
    */
    first_name: string;
    /**
     * @type string
    */
    last_name: string;
    /**
     * @type array
    */
    opportunity_ids: number[];
    /**
     * @type string
    */
    title: string;
};