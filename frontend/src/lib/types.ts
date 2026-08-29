// get mailboxes
export interface Mailbox {
  id: number;
  server_hostport: string;
  username: string;
  primary_inbox: string;
}

export interface Mail {
  id: number;
  mailbox: Mailbox;
  subject: string;
  from_address: string;
  display_name: string; // display name of the sender
  body_parts: Array<Part>;
  recieved_at: string;
}

export interface Part {
  content_type: string; // text/plain, text/html, image/jpeg, etc.
  data?: string; // data (if it's text/plain or text/html)
  link?: string; // link to the content (if it's to a file or image)
}
