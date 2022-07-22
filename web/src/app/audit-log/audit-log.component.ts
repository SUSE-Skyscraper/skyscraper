import { Component, Input, OnChanges, OnInit } from '@angular/core';
import { BackendService, UserItem } from '../backend.service';

@Component({
  selector: 'app-audit-log',
  templateUrl: './audit-log.component.html',
  styleUrls: ['./audit-log.component.scss'],
})
export class AuditLogComponent implements OnInit, OnChanges {
  auditLogs: AuditLogTableItem[] = [];
  auditLogColumns: string[] = [];

  @Input() resource_id?: string;
  @Input() resource_type?: string;

  constructor(private backendService: BackendService) {}

  ngOnInit(): void {
    this.getAuditLogs();
  }

  ngOnChanges() {
    this.getAuditLogs();
  }

  private getAuditLogs() {
    if (this.resource_id !== undefined && this.resource_type !== undefined) {
      this.auditLogColumns = ['caller_id', 'message', 'created_at'];
    } else {
      this.auditLogColumns = [
        'caller_id',
        'resource_type',
        'resource_id',
        'message',
        'created_at',
      ];
    }
    this.backendService
      .getAuditLogs(this.resource_id, this.resource_type)
      .subscribe((response) => {
        if (response.data !== null) {
          let users = new Map<string, UserItem>();

          response.included?.forEach((item) => {
            if (item.type === 'user') {
              users.set(item.id, item as UserItem);
            }
          });

          let tableItems: AuditLogTableItem[] = [];
          response.data.forEach((item) => {
            let caller: string = '';
            switch (item.attributes.caller_type) {
              case 'user':
                caller = `User (${
                  users.get(item.attributes.caller_id)?.attributes.username ??
                  ''
                })`;
                break;
              case 'api_key':
                caller = `API Key (${item.attributes.caller_id})`;
                break;
            }

            tableItems.push({
              message: item.attributes.message,
              caller: caller,
              created_at: new Date(item.attributes.created_at),
              resource_id: item.attributes.resource_id,
              resource_type: item.attributes.resource_type,
            });
          });

          this.auditLogs = tableItems;
        } else {
          this.auditLogs = [];
        }
      });
  }
}

export interface AuditLogTableItem {
  message: string;
  caller: string;
  created_at: Date;
  resource_type: string;
  resource_id: string;
}
