/**
 * Schedules execution of a one-time callback after a delay.
 * 
 * @param fn The function to execute.
 * @param ms Delay in milliseconds.
 * @returns A timer handle that can be passed to clearTimeout.
 */
declare function setTimeout(fn: () => void, ms: number): any;

/**
 * Cancels a timer previously established by setTimeout.
 * @param handle The handle returned by setTimeout.
 */
declare function clearTimeout(handle: any): void;

/**
 * Schedules repeated execution of a callback at a fixed interval.
 * 
 * @param fn The function to execute.
 * @param ms Interval in milliseconds.
 * @returns A timer handle that can be passed to clearInterval.
 */
declare function setInterval(fn: () => void, ms: number): any;

/**
 * Cancels a timer previously established by setInterval.
 * 
 * @param handle The timer handle to cancel.
 */
declare function clearInterval(handle: any): void;
