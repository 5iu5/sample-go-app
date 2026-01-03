export default function TopicCard({title}: TopicCardProps) {
    return(
        <div className="TopicCardDiv">
            <h2>{title}</h2>
        </div>
        
        
    )
}

type TopicCardProps = {
  title: string;
};