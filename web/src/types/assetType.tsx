export type StateMember<TType extends string> = {
  type: TType;
};

export function assertState<TType extends string>(
  state: { type: string },
  ...expectedTypes: TType[]
): asserts state is StateMember<TType> {
  if (!expectedTypes.includes(state.type as TType)) {
    throw new Error(
      `Invalid state ${state.type} (expected one of: ${expectedTypes})`,
    );
  }
}
