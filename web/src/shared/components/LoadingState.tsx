type LoadingStateProps = {
  label: string;
};

export function LoadingState({ label }: LoadingStateProps) {
  return (
    <p className="loading-state" role="status">
      <span className="loading-state__dot" aria-hidden="true" />
      {label}
    </p>
  );
}

import "./LoadingState.css";
