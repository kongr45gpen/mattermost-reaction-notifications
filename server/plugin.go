package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) ReactionHasBeenAdded(c *plugin.Context, post *model.Post, reaction *model.Reaction) {
	user, _ := p.API.GetUser(reaction.UserId)
	channel, _ := p.API.GetChannel(post.ChannelId)
	team, _ := p.API.GetTeam(channel.TeamId)
	directChannel, _ := p.API.GetDirectChannel(post.UserId, reaction.UserId)

	config := p.API.GetConfig()
	p.API.LogDebug(fmt.Sprintf("Config server URL %s", *config.ServiceSettings.SiteURL))

	link := fmt.Sprintf("%s/%s/pl/%s", *config.ServiceSettings.SiteURL, team.Name, post.Id)

	p.API.SendEphemeralPost(post.UserId, &model.Post{
		UserId:    reaction.UserId,
		ChannelId: post.ChannelId,
		Message:   fmt.Sprintf("@%s has reacted with :%s: to [your message](%s)", user.Username, reaction.EmojiName, link),
	})

	_, err := p.API.CreatePost(&model.Post{
		UserId:    reaction.UserId,
		ChannelId: directChannel.Id,
		Message:   fmt.Sprintf("@%s has reacted with :%s: to [your message](%s)", user.Username, reaction.EmojiName, link),
	})

	if err != nil {
		p.API.LogError(err.DetailedError)
	}

	// p.API.DeletePost(newPost.Id)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
