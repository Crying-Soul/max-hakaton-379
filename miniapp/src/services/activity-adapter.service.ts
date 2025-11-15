// src/services/activity-adapter.service.ts
import { Activity } from '../common/types/domain';
import { MapEvent } from '../common/types/api';

export class ActivityAdapterService {
  static mapEventToActivity(event: MapEvent): Activity {
    return {
      id: event.id.toString(), // Конвертируем number в string для фронта
      title: event.title,
      description: event.description,
      date: event.date, // Используем date как есть
      durationHours: event.durationHours,
      location: event.location,
      lat: event.locationLat,
      lon: event.locationLon,
      status: event.status === 'open' && event.slotsLeft > 0 ? 'open' : 'closed',
      distanceText: `${event.distanceKm.toFixed(1)} км`,
      distanceKm: event.distanceKm,
      categoryId: event.categoryId,
      categoryName: event.categoryName,
      organizerId: event.organizerId,
      contacts: event.contacts,
      maxVolunteers: event.maxVolunteers,
      currentVolunteers: event.currentVolunteers,
      slotsLeft: event.slotsLeft,
      deeplink: this.generateDeeplink(event.id),
      createdAt: event.createdAt,
      updatedAt: event.updatedAt,
    };
  }

  private static generateDeeplink(eventId: number): string {
    // Генерируем deeplink на основе ID события
    return `https://max.ru/t379_hakaton_bot?start=OpenEvent_${eventId}`;
  }

  static mapEventsToActivities(events: MapEvent[]): Activity[] {
    return events.map(event => this.mapEventToActivity(event));
  }
}