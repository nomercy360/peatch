export type CreateCollaboration = {
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
     * @type boolean | undefined
    */
    is_payable?: boolean;
    /**
     * @type integer
    */
    opportunity_id: number;
    /**
     * @type string
    */
    title: string;
};