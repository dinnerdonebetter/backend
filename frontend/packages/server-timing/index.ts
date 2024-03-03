export const ServerTimingHeaderName = 'Server-Timing';

export class ServerTimingEvent {
  name: string = '';
  description: string = '';
  startTime: Date = new Date();
  endTime: Date = new Date();
  duration: number = 0;

  constructor(name: string = '', description: string = '') {
    this.name = name;
    this.description = description;
    this.startTime = new Date();
  }

  end() {
    this.endTime = new Date();
    this.duration = this.endTime.getTime() - this.startTime.getTime();
  }
}

export class ServerTiming {
  events: ServerTimingEvent[] = [];

  addEvent(name: string = '', description: string = ''): ServerTimingEvent {
    const event = new ServerTimingEvent(name, description);
    this.events.push(event);
    return event;
  }

  headerValue(): string {
    let result = '';
    for (const event of this.events) {
      result += `${event.name};desc="${event.description}";dur=${event.duration},`;
    }
    return result.trim().replace(/,(\s*)$/, '$1');
  }
}
