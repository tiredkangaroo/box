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
  body: string;
  recieved_at: string;
}
